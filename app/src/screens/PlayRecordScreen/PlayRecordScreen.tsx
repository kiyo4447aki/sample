import colors from '@/constants/colors';
import { useAuth } from '@/contexts/AuthContext';
import { RootStackParamList } from '@/navigation/RootNavigator';
import { RouteProp, useRoute } from '@react-navigation/native';
import { lockAsync, OrientationLock } from 'expo-screen-orientation';
import React, { useEffect } from 'react';
import { StyleSheet, View } from 'react-native';
import { ResizeMode } from 'react-native-video';
import VideoPlayer from 'react-native-video-controls';

const PlayRecordScreen = () => {
  const route = useRoute<RouteProp<RootStackParamList, 'PlayRecord'>>();
  const { token } = useAuth();
  useEffect(() => {
    lockAsync(OrientationLock.LANDSCAPE);

    return () => {
      lockAsync(OrientationLock.PORTRAIT_UP);
    };
  }, []);

  return (
    <View style={styles.container}>
      <VideoPlayer
        source={{
          uri: route.params.url,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }}
        disableFullscreen={true}
        disableVolume={true}
        controlAnimationTiming={300}
        controlTimeout={5000}
        onError={() => {}}
        resizeMode={ResizeMode.CONTAIN}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.dark.header.statusbar,
    flex: 1,
  },
});

export default PlayRecordScreen;
