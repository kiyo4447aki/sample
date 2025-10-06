import colors from '@/constants/colors';
import { useAuth } from '@/contexts/AuthContext';
import { useDevices } from '@/contexts/DevicesContext';
import DatePicker from '@/screens/RecListScreen/components/DatePicker/DatePicker';
import Dropdown from '@/screens/RecListScreen/components/Dropdown/Dropdown';
import { SearchContext } from '@/screens/RecListScreen/contexts/SearchContext';
import initclient, { RecordListResponse } from '@/services/recorderApi';
import { useContext } from 'react';
import { Alert, StyleSheet, Text, TouchableOpacity, View } from 'react-native';

const SearchCard = () => {
  const ctx = useContext(SearchContext);
  if (!ctx) {
    throw new Error('Search contextはProviderの中で使用してください。');
  }
  const { setSelectedCamera, setDate, setRecords } = ctx;
  const { device } = useDevices();
  const { token } = useAuth();

  const client = initclient(device, token);

  const clearParameters = () => {
    setSelectedCamera('');
    setDate(new Date());
  };

  const fetchRecordList = async () => {
    try {
      const res = await client.get<RecordListResponse>('/api/records');
      setRecords(res.data.records);
    } catch {
      Alert.alert(
        '通信エラー',
        'データの取得に失敗しました',
        [
          {
            text: 'OK',
          },
        ],
        { cancelable: false },
      );
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>検索条件</Text>
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionTitle}>日付</Text>
        <DatePicker />
      </View>
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionTitle}>カメラ</Text>
        <Dropdown />
      </View>
      <View style={styles.buttonsContainer}>
        <TouchableOpacity style={styles.button} onPress={clearParameters}>
          <Text style={styles.buttonText}>条件をクリア</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.buttonBk} onPress={fetchRecordList}>
          <Text style={styles.buttonBkText}>検索</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  button: {
    alignItems: 'center',
    alignSelf: 'stretch',
    backgroundColor: colors.light.card.bg,
    borderColor: colors.light.card.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
    justifyContent: 'center',
    paddingVertical: 8,
  },
  buttonBk: {
    alignItems: 'center',
    alignSelf: 'stretch',
    backgroundColor: colors.light.buttonBk.bg,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
    justifyContent: 'center',
    paddingVertical: 8,
  },
  buttonBkText: {
    color: colors.light.buttonBk.text,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
  buttonText: {
    color: colors.light.button.text,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
  buttonsContainer: {
    gap: 10,
  },
  container: {
    backgroundColor: colors.light.card.bg,
    borderColor: colors.light.card.border,
    borderRadius: 14,
    borderStyle: 'solid',
    borderWidth: 1,
    gap: 10,
    padding: 20,
  },
  sectionContainer: {
    gap: 5,
  },
  sectionTitle: {
    color: colors.light.card.sectionTitle,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 28,
  },
  title: {
    color: colors.light.card.primaryText,
    fontSize: 18,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 28,
  },
});

export default SearchCard;
