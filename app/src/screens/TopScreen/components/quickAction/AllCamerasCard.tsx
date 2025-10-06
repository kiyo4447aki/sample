import tableIcon from '@/../assets/icons/table.png';
import Card from '@/screens/TopScreen/components/Card/Card';
import React from 'react';
import { Dimensions, StyleSheet, View } from 'react-native';

const AllCamerasCard = () => {
  return (
    <Card>
      <View style={styles.leftItemsContainer}>
        <Card.Icon source={tableIcon} />
        <Card.TextContainer>
          <Card.Title text="全カメラ同時視聴" />
          <Card.Subtitle text="すべてのカメラを同時に表示" />
        </Card.TextContainer>
      </View>
      <Card.Arrow />
    </Card>
  );
};

const styles = StyleSheet.create({
  leftItemsContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: Dimensions.get('window')['width'] / 2 - 130,
  },
});

export default AllCamerasCard;
