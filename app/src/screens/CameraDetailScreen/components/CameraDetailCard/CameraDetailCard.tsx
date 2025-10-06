import colors from '@/constants/colors';
import React from 'react';
import { StyleSheet, Text, View } from 'react-native';

const CameraDetailCard = () => {
  return (
    <View style={styles.container}>
      <View style={styles.sectionContainer}>
        <Text style={styles.label}>カメラID</Text>
        <Text style={styles.value}>cam-01</Text>
      </View>
      <View style={styles.sectionContainer}>
        <Text style={styles.label}>解像度</Text>
        <Text style={styles.value}>1920×1080</Text>
      </View>
      <View style={styles.sectionContainer}>
        <Text style={styles.label}>フレームレート</Text>
        <Text style={styles.value}>30FPS</Text>
      </View>
      <View style={styles.sectionContainer}>
        <Text style={styles.label}>接続状況</Text>
        <Text style={styles.connectionValue}>オンライン</Text>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  connectionValue: {
    color: colors.dark.connectionIndicator.text,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: '300',
    lineHeight: 20,
  },
  container: {
    alignItems: 'flex-start',
    alignSelf: 'stretch',
    backgroundColor: colors.dark.card.bg,
    borderColor: colors.dark.card.border,
    borderRadius: 14,
    borderStyle: 'solid',
    borderWidth: 1,
    flexDirection: 'row',
    flexWrap: 'wrap',
    padding: 18,
  },
  label: {
    color: colors.dark.card.secondaryText,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: '300',
    lineHeight: 20,
  },
  sectionContainer: {
    alignItems: 'flex-start',
    justifyContent: 'flex-start',
    width: '50%',
  },
  value: {
    color: colors.dark.card.primaryText,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: '400',
    lineHeight: 20,
  },
});

export default CameraDetailCard;
