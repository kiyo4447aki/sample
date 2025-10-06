import colors from '@/constants/colors';
import React from 'react';
import { StyleSheet, Text } from 'react-native';

const Title = ({ text }: { text: string }) => {
  return <Text style={styles.text}>{text}</Text>;
};

const styles = StyleSheet.create({
  text: {
    color: colors.light.card.sectionTitle,
    fontSize: 18,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 28,
  },
});

export default Title;
