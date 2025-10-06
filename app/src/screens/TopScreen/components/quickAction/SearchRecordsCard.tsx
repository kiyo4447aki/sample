import videoIcon from '@/../assets/icons/video-green.png';
import { RootStackParamList } from '@/navigation/RootNavigator';
import Card from '@/screens/TopScreen/components/Card/Card';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import React from 'react';
import { Dimensions, StyleSheet, View } from 'react-native';

type NavProp = StackNavigationProp<RootStackParamList, 'Top'>;

const SearchRecordsCard = () => {
  const navigation = useNavigation<NavProp>();
  return (
    <Card
      onPress={() => {
        navigation.navigate('RecList');
      }}
    >
      <View style={styles.leftItemsContainer}>
        <Card.Icon source={videoIcon} />
        <Card.TextContainer>
          <Card.Title text="録画映像検索" />
          <Card.Subtitle text="過去の録画を検索・視聴　　"></Card.Subtitle>
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

export default SearchRecordsCard;
