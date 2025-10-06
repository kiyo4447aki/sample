import logoIcon from '@/../assets/icons/logo.png';
import colors from '@/constants/colors';
import React from 'react';
import { Image, StyleSheet, Text, View } from 'react-native';

const Header = ({ deviceId }: { deviceId?: string }) => {
  return (
    <View style={styles.container}>
      <Image source={logoIcon} style={styles.logo} />
      {deviceId && (
        <View style={styles.deviceIdContainer}>
          <Text style={styles.deviceIdText}>デバイスID：{deviceId}</Text>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    backgroundColor: colors.light.header.bg,
    flexDirection: 'row',
    height: 56,
    justifyContent: 'space-between',
    paddingHorizontal: 10,
  },
  deviceIdContainer: {
    alignItems: 'center',
    borderColor: colors.light.header.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
    justifyContent: 'center',
    paddingHorizontal: 9,
    paddingVertical: 3,
  },
  deviceIdText: {
    color: colors.light.header.secondaryText,
    fontSize: 11,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 16,
    textAlign: 'center',
  },
  logo: {
    height: 24,
    width: 100,
  },
});

export default Header;
