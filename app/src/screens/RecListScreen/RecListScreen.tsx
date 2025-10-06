import Title from '@/components/Title/Title';
import colors from '@/constants/colors';
import { DropdownItem } from '@/screens/RecListScreen/components/Dropdown/Dropdown';
import Header from '@/screens/RecListScreen/components/Header/Header';
import SearchCard from '@/screens/RecListScreen/components/SearchCard/SearchCard';
import { SearchContext } from '@/screens/RecListScreen/contexts/SearchContext';
import { Record } from '@/services/recorderApi';
import React, { useState } from 'react';
import { StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

const RecListScreen = () => {
  const [selectedCamera, setSelectedCamera] = useState<string>('');
  const [cameras, setCameras] = useState<DropdownItem[]>([
    { label: 'カメラ1', value: 'cam-01' },
    { label: 'カメラ2', value: 'cam-02' },
  ]);
  const [date, setDate] = useState(new Date());
  const [records, setRecords] = useState<Record[]>([]);

  return (
    <SearchContext.Provider
      value={{
        selectedCamera,
        setSelectedCamera,
        cameras,
        setCameras,
        date,
        setDate,
        records,
        setRecords,
      }}
    >
      <View style={styles.container}>
        <SafeAreaView edges={['top']} style={styles.safeAreaTop} />
        <Header />
        <View style={styles.contentArea}>
          <SearchCard />
          <Title text="検索結果" />
        </View>
        <SafeAreaView edges={['bottom']} style={styles.safeAreaBottom} />
      </View>
    </SearchContext.Provider>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.light.header.bg,
    flex: 1,
    justifyContent: 'flex-start',
  },
  contentArea: {
    backgroundColor: colors.light.bg,
    flex: 1,
    gap: 20,
    paddingHorizontal: 10,
    paddingVertical: 10,
  },
  safeAreaBottom: {
    backgroundColor: colors.light.bg,
  },
  safeAreaTop: {
    backgroundColor: colors.light.header.bg,
  },
});

export default RecListScreen;
