import React from 'react';
import { Image, StyleSheet, TouchableOpacity } from 'react-native';

import allowIcon from '@/../assets/icons/allow-left-wh.png';
import colors from '@/constants/colors';
import { useNavigation } from '@react-navigation/native';

const BackButton = () => {
  const navigation = useNavigation();

  return (
    <TouchableOpacity style={styles.container} onPress={() => navigation.goBack()}>
      <Image source={allowIcon} style={styles.icon} />
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.dark.button.bgSolid,
    borderRadius: 20,
    left: 16,
    padding: 8,
    position: 'absolute',
    top: 20,
    zIndex: 10,
  },
  icon: {
    height: 15,
    width: 15,
  },
});

export default BackButton;
