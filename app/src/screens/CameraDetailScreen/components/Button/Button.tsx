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
    <TouchableOpacity style={styles.button} onPress={onPress}>
      <Image source={icon} style={styles.icon} />
      <Text style={styles.text}>{text}</Text>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    alignItems: 'center',
    backgroundColor: colors.dark.button.bg,
    borderColor: colors.dark.button.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
    flexDirection: 'row',
    gap: 12,
    height: 50,
    justifyContent: 'center',
    width: 160,
  },
  icon: {
    height: 16,
    width: 16,
  },
  text: {
    color: colors.dark.button.text,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
    textAlign: 'center',
  },
});

export default Button;
