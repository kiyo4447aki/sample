import arrow from '@/../assets/icons/allow-left-bk.png';
import colors from '@/constants/colors';
import { useNavigation } from '@react-navigation/native';
import React from 'react';
import { Image, StyleSheet, Text, TouchableOpacity, View } from 'react-native';

const Header = () => {
  const navigation = useNavigation();
  return (
    <View style={styles.container}>
      <TouchableOpacity
        onPress={() => {
          navigation.goBack();
        }}
      >
        <Image source={arrow} style={styles.arrow} />
      </TouchableOpacity>
      <Text style={styles.text}>録画映像検索</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  arrow: {
    height: 22,
    width: 22,
  },
  container: {
    alignItems: 'center',
    backgroundColor: colors.light.header.bg,
    flexDirection: 'row',
    gap: 12,
    height: 56,
    justifyContent: 'flex-start',
    paddingHorizontal: 10,
  },
  text: {
    color: colors.light.header.primaryText,
    fontSize: 20,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 28,
  },
});

export default Header;
