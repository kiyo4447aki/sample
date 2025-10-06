//TODOデータ取得との連携

import Title from '@/components/Title/Title';
import colors from '@/constants/colors';
import { useAuth } from '@/contexts/AuthContext';
import { useDevices } from '@/contexts/DevicesContext';
import CameraCard from '@/screens/TopScreen/components/CameraCard/CameraCard';
import Header from '@/screens/TopScreen/components/Header/Header';
import AllCamerasCard from '@/screens/TopScreen/components/quickAction/AllCamerasCard';
import SearchRecordsCard from '@/screens/TopScreen/components/quickAction/SearchRecordsCard';
import React, { useEffect } from 'react';
import { Alert, DevSettings, StyleSheet, View } from 'react-native';
import { ScrollView } from 'react-native-gesture-handler';
import { SafeAreaView } from 'react-native-safe-area-context';

const TopScreen = () => {
  const { isLoggedIn } = useAuth();
  const { device, fetchDeviceList } = useDevices();
  useEffect(() => {
    if (isLoggedIn) {
      try {
        fetchDeviceList();
      } catch {
        Alert.alert(
          '通信エラー',
          'データの取得に失敗しました',
          [
            {
              text: 'OK',
              onPress: () => {
                DevSettings.reload();
              },
            },
          ],
          { cancelable: false },
        );
      }
    }
  }, [isLoggedIn]);

  useEffect(() => {
    console.log(device);
  }, [device]);

  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeAreaTop} />
      <Header deviceId={device} />
      <View style={styles.contentArea}>
        <View style={styles.quickActionContainer}>
          <Title text="クイックアクション" />
          <AllCamerasCard />
          <SearchRecordsCard />
        </View>
        <View style={styles.cameraListContainer}>
          <Title text="カメラ一覧" />
          <ScrollView style={styles.scroll} contentContainerStyle={styles.cameraList}>
            <CameraCard cameraId="cam-01" location="" />
            <CameraCard cameraId="cam-02" location="" />
            <CameraCard cameraId="cam-03" location="" />
            <CameraCard cameraId="cam-04" location="" />
            <CameraCard cameraId="cam-05" location="" />
            <CameraCard cameraId="cam-06" location="" />
          </ScrollView>
        </View>
      </View>
      <SafeAreaView edges={['bottom']} style={styles.safeAreaBottom} />
    </View>
  );
};

const styles = StyleSheet.create({
  cameraList: {
    gap: 4,
  },
  cameraListContainer: {
    flex: 1,
    gap: 10,
  },
  container: {
    flex: 1,
    justifyContent: 'flex-start',
  },
  contentArea: {
    backgroundColor: colors.light.bg,
    flex: 1,
    gap: 15,
    paddingHorizontal: 10,
    paddingVertical: 10,
  },
  quickActionContainer: {
    gap: 4,
  },
  safeAreaBottom: {
    backgroundColor: colors.light.bg,
  },
  safeAreaTop: {
    backgroundColor: colors.light.header.bg,
  },
  scroll: {
    flex: 1,
  },
});

export default TopScreen;
