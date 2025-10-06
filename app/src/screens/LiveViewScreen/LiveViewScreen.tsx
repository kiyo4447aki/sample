import BackButton from '@/components/BackButton/BackButton';
import JanusViewer from '@/components/JanusViewer/JanusViewer';
import colors from '@/constants/colors';
import { ConnectionInfoType, useDevices } from '@/contexts/DevicesContext';
import { RootStackParamList } from '@/navigation/RootNavigator';
import { RouteProp, useNavigation, useRoute } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { lockAsync, OrientationLock } from 'expo-screen-orientation';
import React, { useEffect, useState } from 'react';
import { Alert, StyleSheet, View } from 'react-native';

type NavProp = StackNavigationProp<RootStackParamList, 'CameraDetail'>;

const LiveViewScreen = () => {
  const [connectionInfo, setConnectionInfo] = useState<ConnectionInfoType | undefined>();
  const route = useRoute<RouteProp<RootStackParamList, 'CameraDetail'>>();
  const navigation = useNavigation<NavProp>();
  const { getConnectionInfo } = useDevices();

  const cameraId = route.params.cameraId;

  useEffect(() => {
    lockAsync(OrientationLock.LANDSCAPE);

    (async () => {
      try {
        const connInfo = await getConnectionInfo();
        setConnectionInfo(connInfo);
      } catch {
        Alert.alert(
          '通信エラー',
          'データの取得に失敗しました',
          [
            {
              text: 'OK',
              onPress: () => {
                navigation.goBack();
              },
            },
          ],
          { cancelable: false },
        );
      }
    })();

    return () => {
      lockAsync(OrientationLock.PORTRAIT_UP);
    };
  }, []);

  return (
    <View style={styles.container}>
      <BackButton />
      {connectionInfo ? (
        <JanusViewer
          url={connectionInfo.janusUrl}
          roomId={connectionInfo.deviceId}
          password={connectionInfo.password}
          feedId={cameraId}
        />
      ) : undefined}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.dark.header.statusbar,
    flex: 1,
  },
});

export default LiveViewScreen;
