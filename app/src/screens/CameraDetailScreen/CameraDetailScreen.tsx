import fullscreenIcon from '@/../assets/icons/fullscreen.png';
import muteIcon from '@/../assets/icons/mute.png';
import reloadIcon from '@/../assets/icons/reload.png';
import JanusViewer from '@/components/JanusViewer/JanusViewer';
import colors from '@/constants/colors';
import { ConnectionInfoType, useDevices } from '@/contexts/DevicesContext';
import { RootStackParamList } from '@/navigation/RootNavigator';

import Button from '@/screens/CameraDetailScreen/components/Button/Button';
import CameraDetailCard from '@/screens/CameraDetailScreen/components/CameraDetailCard/CameraDetailCard';
import Header from '@/screens/CameraDetailScreen/components/Header/Header';
import { RouteProp, useNavigation, useRoute } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import React, { useEffect, useState } from 'react';
import { Alert, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

type NavProp = StackNavigationProp<RootStackParamList, 'CameraDetail'>;

const CameraDetailScreen = () => {
  const [connectionInfo, setConnectionInfo] = useState<ConnectionInfoType | undefined>();
  const route = useRoute<RouteProp<RootStackParamList, 'CameraDetail'>>();
  const navigation = useNavigation<NavProp>();
  const { getConnectionInfo } = useDevices();

  const cameraId = route.params.cameraId;

  const fetchConnInfo = async () => {
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
  };

  useEffect(() => {
    fetchConnInfo();
  }, []);

  return (
    <SafeAreaView style={styles.container}>
      <Header />
      <View style={styles.contentArea}>
        <View style={styles.videoArea}>
          {connectionInfo ? (
            <JanusViewer
              url={connectionInfo.janusUrl}
              roomId={connectionInfo.deviceId}
              password={connectionInfo.password}
              feedId={cameraId}
            />
          ) : undefined}
        </View>
        <Button
          icon={fullscreenIcon}
          text="フルスクリーン"
          onPress={() => navigation.navigate('LiveView', { cameraId })}
        />
        <View style={styles.bottomButonsContainer}>
          <Button icon={muteIcon} text="ミュート解除" />
          <Button
            icon={reloadIcon}
            text="更新"
            onPress={() => {
              setConnectionInfo(undefined);
              fetchConnInfo();
            }}
          />
        </View>
        <CameraDetailCard />
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  bottomButonsContainer: {
    alignItems: 'center',
    alignSelf: 'stretch',
    flexDirection: 'row',
    justifyContent: 'space-around',
  },

  container: {
    backgroundColor: colors.dark.header.statusbar,
    flex: 1,
    justifyContent: 'flex-start',
  },
  contentArea: {
    alignItems: 'center',
    flex: 1,
    gap: 16,
    padding: 16,
  },
  /* eslint-disable react-native/no-color-literals */
  videoArea: {
    alignSelf: 'stretch',
    backgroundColor: '#0F172B',
    borderColor: '#314158',
    borderRadius: 14,
    borderStyle: 'solid',
    borderWidth: 1,
    height: 210,
  },
  /* eslint-enable react-native/no-color-literals */
});

export default CameraDetailScreen;
