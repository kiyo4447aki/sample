import arrow from '@/../assets/icons/allow-left-wh.png';
import connIcon from '@/../assets/icons/conn-good.png';
import colors from '@/constants/colors';
import { useNavigation } from '@react-navigation/native';
import React from 'react';
import { Image, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
const Header = () => {
  const navigation = useNavigation();

  return (
    <View style={styles.container}>
      <View style={styles.leftItemContainer}>
        <TouchableOpacity
          onPress={() => {
            navigation.goBack();
          }}
        >
          <Image source={arrow} style={styles.arrow} />
        </TouchableOpacity>
        <View>
          <Text style={styles.title}>cam-01</Text>
          <Text style={styles.subTitle}>エントランス</Text>
        </View>
      </View>
      <View style={styles.indicator}>
        <Image source={connIcon} style={styles.indicatorIcon} />
        <Text style={styles.indicatorText}>ライブ</Text>
      </View>
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
    alignSelf: 'stretch',
    backgroundColor: colors.dark.header.bg,
    borderBottomColor: colors.dark.header.border,
    borderBottomWidth: 1,
    borderStyle: 'solid',
    flexDirection: 'row',
    gap: 12,
    height: 56,
    justifyContent: 'space-between',
    paddingHorizontal: 10,
  },
  indicator: {
    alignItems: 'center',
    backgroundColor: colors.dark.connectionIndicator.bg,
    borderRadius: 33554400,
    flexDirection: 'row',
    gap: 4,
    paddingHorizontal: 8,
    paddingVertical: 4,
  },
  indicatorIcon: {
    height: 12,
    width: 12,
  },
  indicatorText: {
    color: colors.dark.connectionIndicator.text,
    fontSize: 12,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 16,
  },
  leftItemContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: 12,
  },
  subTitle: {
    color: colors.dark.header.secondaryText,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
  title: {
    color: colors.dark.header.primaryText,
    fontSize: 18,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 28,
  },
});

export default Header;
