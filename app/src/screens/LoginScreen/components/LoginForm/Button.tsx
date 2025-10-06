import colors from '@/constants/colors';
import React from 'react';
import { Image, ImageSourcePropType, StyleSheet, Text, TouchableOpacity } from 'react-native';

type ButtonProps = {
  icon: ImageSourcePropType;
  text: string;
  onPress?: () => void;
};

const Button = ({ icon, text, onPress }: ButtonProps) => {
  return (
    <TouchableOpacity style={styles.container} onPress={onPress}>
      <Image source={icon} style={styles.icon} />
      <Text style={styles.text}>{text}</Text>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    alignSelf: 'stretch',
    backgroundColor: colors.light.buttonBlue.bg,
    borderRadius: 8,
    flexDirection: 'row',
    gap: 8,
    height: 36,
    justifyContent: 'center',
    paddingHorizontal: 16,
    paddingVertical: 8,
  },
  icon: {
    height: 16,
    width: 16,
  },
  text: {
    color: colors.light.buttonBlue.text,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
    textAlign: 'center',
  },
});

export default Button;
