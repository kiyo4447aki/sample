import blackCalenderIcon from '@/../assets/icons/calender-black.png';
import colors from '@/constants/colors';
import { SearchContext } from '@/screens/RecListScreen/contexts/SearchContext';
import DateTimePicker, { DateTimePickerEvent } from '@react-native-community/datetimepicker';
import React, { useContext, useState } from 'react';
import { Image, Platform, StyleSheet, Text, TouchableOpacity } from 'react-native';

const DatePicker = () => {
  const [isVisible, setIsVisible] = useState(false);
  const ctx = useContext(SearchContext);
  if (!ctx) {
    throw new Error('Search contextはProviderの中で使用してください。');
  }
  const { date, setDate } = ctx;

  const onChange = (event: DateTimePickerEvent, selectedDate?: Date) => {
    setIsVisible(false);
    if (selectedDate) {
      setDate(selectedDate);
    }
  };

  return (
    <>
      {isVisible ? (
        <DateTimePicker
          value={date}
          mode="date"
          display={Platform.OS === 'ios' ? 'spinner' : 'default'}
          onChange={onChange}
          locale="ja-JP"
        />
      ) : (
        <TouchableOpacity
          style={styles.button}
          onPress={() => {
            setIsVisible(true);
          }}
        >
          <Text style={styles.text}>{date.toLocaleDateString('ja-jp')}</Text>
          <Image source={blackCalenderIcon} style={styles.icon} />
        </TouchableOpacity>
      )}
    </>
  );
};

const styles = StyleSheet.create({
  button: {
    alignItems: 'center',
    alignSelf: 'stretch',
    backgroundColor: colors.light.button.bg,
    borderColor: colors.light.button.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
    flexDirection: 'row',
    gap: 4,
    paddingHorizontal: 15,
    paddingVertical: 10,
  },
  icon: {
    height: 16,
    width: 16,
  },
  text: {
    color: colors.light.button.text,
    fontSize: 16,
    fontStyle: 'normal',
    fontWeight: 400,
    lineHeight: 24,
  },
});

export default DatePicker;
