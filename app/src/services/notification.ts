import api from '@/services/api';
import { isString } from '@/services/common';
import notifee, { AndroidImportance, AuthorizationStatus, EventType } from '@notifee/react-native';
import messaging from '@react-native-firebase/messaging';
import { Alert, Linking, Platform } from 'react-native';

const ensureNotifyPermission = async (): Promise<boolean> => {
  let settings = await notifee.getNotificationSettings();
  if (settings.authorizationStatus < AuthorizationStatus.AUTHORIZED) {
    settings = await notifee.requestPermission();
  }
  return settings.authorizationStatus >= AuthorizationStatus.AUTHORIZED;
};

const createAndroidChannel = async () => {
  if (Platform.OS !== 'android') return;
  await notifee.createChannel({
    id: 'alerts',
    name: 'Alerts',
    importance: AndroidImportance.HIGH,
    vibration: true,
  });
};

const registerTokenToBackend = async (token: string) => {
  try {
    await api.post('/notify/token/register', {
      token,
      platform: Platform.OS,
    });
  } catch (e) {
    if (e instanceof Error) {
      throw new Error('プッシュトークンの登録に失敗しました：' + e.message);
    } else {
      throw new Error('不明なエラーによりプッシュトークンの登録に失敗しました');
    }
  }
};

//アプリ起動直後に呼び出し
export const initPush = async () => {
  const isPermissionGranted = await ensureNotifyPermission();
  if (!isPermissionGranted) {
    Alert.alert('通知が無効です', '端末の設定から通知を有効にしてください');
  }
  await messaging().registerDeviceForRemoteMessages();
  await createAndroidChannel();

  //トークンの取得と登録
  const token = await messaging().getToken();
  try {
    await registerTokenToBackend(token);
  } catch {
    Alert.alert('通知の購読に失敗しました', 'アプリを再起動してください');
  }

  //トークン変更時
  const unsubRefresh = messaging().onTokenRefresh(async newToken => {
    try {
      await registerTokenToBackend(newToken);
    } catch {
      Alert.alert('通知の購読の更新に失敗しました', 'アプリを再起動してください');
    }
  });

  const unsubOnMessage = messaging().onMessage(async msg => {
    const title = msg.notification?.title ?? '';
    const body = msg.notification?.body ?? '';
    await notifee.displayNotification({
      title,
      body,
      data: msg.data ?? {},
      android: {
        channelId: 'alerts',
        pressAction: { id: 'default' },
      },
    });
  });

  //TODO 通知の受信履歴を保存

  // Notifeeの通知タップ（ローカル通知）を拾う
  const unsubNotifee = notifee.onForegroundEvent(({ type, detail }) => {
    if (type === EventType.PRESS) {
      const path = getDeepLinkPath(detail.notification?.data);
      openDeepLink(path);
    }
  });
  //バックグラウンド起動時
  const unsubOpened = messaging().onNotificationOpenedApp(msg => {
    const path = getDeepLinkPath(msg.data);
    openDeepLink(path);
  });

  //アプリを通知から起動したかどうか判定
  const initial = await messaging().getInitialNotification();
  if (initial) {
    const path = getDeepLinkPath(initial.data);
    openDeepLink(path);
  }

  const unsubscribe = () => {
    unsubRefresh();
    unsubOnMessage();
    unsubOpened();
    unsubNotifee();
  };
  console.log(token);

  return unsubscribe;
};

const getDeepLinkPath = (data?: { [k: string]: string | object | number }) => {
  let path: string | undefined = undefined;
  const deviceId = data?.deviceId;
  const eventId = data?.eventId;
  if (deviceId && eventId && isString(deviceId) && isString(eventId)) {
    path = `/devices/${deviceId}/alerts/${eventId}`;
  } else if (deviceId && isString(deviceId)) {
    path = `/devices/${deviceId}`;
  } else {
    return undefined;
  }

  return `myapp://${path}`;
};

const openDeepLink = async (path: string | undefined) => {
  if (!path) return;
  let isValidLink: boolean = false;
  isValidLink = await Linking.canOpenURL(path);
  if (isValidLink) await Linking.openURL(path);
};
