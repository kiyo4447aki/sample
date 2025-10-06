import colors from '@/constants/colors';
import { SearchContext } from '@/screens/RecListScreen/contexts/SearchContext';
import React, { useContext, useState } from 'react';
import { StyleSheet, ViewStyle } from 'react-native';
import DropDownPicker from 'react-native-dropdown-picker';

export type DropdownItem = {
  label: string;
  value: string;
};

const Dropdown = () => {
  const [isOpen, setIsOpen] = useState(false);
  const ctx = useContext(SearchContext);
  if (!ctx) {
    throw new Error('Search contextはProviderの中で使用してください。');
  }
  const { selectedCamera, setSelectedCamera, cameras, setCameras } = ctx;

  return (
    <DropDownPicker
      open={isOpen}
      value={selectedCamera}
      items={cameras}
      setOpen={setIsOpen}
      setValue={setSelectedCamera}
      setItems={setCameras}
      zIndex={1000}
      zIndexInverse={3000}
      placeholder="カメラを選択"
      style={styles.dropdown}
      textStyle={styles.text}
      dropDownContainerStyle={styles.dropdownContainer}
      //ライブラリの型定義が間違えているためアサーション
      arrowIconStyle={styles.icon as ViewStyle}
    />
  );
};

const styles = StyleSheet.create({
  dropdown: {
    backgroundColor: colors.light.card.bg,
    borderColor: colors.light.button.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
  },
  dropdownContainer: {
    backgroundColor: colors.light.button.bg,
    borderColor: colors.light.button.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
  },
  icon: {
    tintColor: colors.light.other.dropdownIcon,
  },
  text: {
    color: colors.light.other.dropdownIcon,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
});

export default Dropdown;
