import Card from '@/screens/TopScreen/components/Card/Card';
import React from 'react';

import cameraIcon from '@/../assets/icons/camera.png';
import { RootStackParamList } from '@/navigation/RootNavigator';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { Dimensions, StyleSheet, View } from 'react-native';

type NavProp = StackNavigationProp<RootStackParamList, 'Top'>;

const CameraCard = ({ cameraId, location }: { cameraId: string; location: string }) => {
  const navigation = useNavigation<NavProp>();

  return (
    <Card
      style={styles.card}
      onPress={() =>
        navigation.navigate('CameraDetail', {
          cameraId,
        })
      }
    >
      <View style={styles.leftItemsContainer}>
        <Card.Icon source={cameraIcon} />
        <Card.TextContainer>
          <Card.Title text={cameraId} />
          <Card.Subtitle text={location} />
        </Card.TextContainer>
      </View>
      <Card.Arrow />
    </Card>
  );
};

const styles = StyleSheet.create({
  card: {
    paddingVertical: 12,
  },
  leftItemsContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: Dimensions.get('window')['width'] / 2 - 95,
  },
});

export default CameraCard;
