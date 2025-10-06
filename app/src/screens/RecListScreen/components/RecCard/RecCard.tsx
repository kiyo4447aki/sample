import grayCalenderIcon from '@/../assets/icons/calender-gray.png';
import downloadIcon from '@/../assets/icons/download.png';
import playIcon from '@/../assets/icons/play.png';
import colors from '@/constants/colors';
import { Record } from '@/services/recorderApi';
import React from 'react';
import { Image, StyleSheet, Text, TouchableOpacity, View } from 'react-native';

type RecCardProps = {
  record: Record;
};

const RecCard = ({ record }: RecCardProps) => {
  return (
    <View style={styles.container}>
      <View style={styles.primaryTextsContainer}>
        <Text style={styles.cameraName}>{record.name}</Text>
        <View style={styles.dateTimeContainer}>
          <Image source={grayCalenderIcon} style={styles.calenderIcon} />
          <Text style={styles.dateTime}>2024-01-15 14:30:25</Text>
        </View>
      </View>
      <View style={styles.secondaryTextsContainer}>
        <Text style={styles.detailText}>時間：00:05:42</Text>
        <Text style={styles.detailText}>サイズ：125MB</Text>
      </View>
      <View style={styles.buttonsContainer}>
        <TouchableOpacity style={styles.playButton}>
          <Image source={playIcon} style={styles.buttonIcon} />
          <Text style={styles.playButtonText}>再生</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.downloadButton}>
          <Image source={downloadIcon} style={styles.buttonIcon} />
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  buttonIcon: {
    height: 16,
    width: 16,
  },
  buttonsContainer: {
    alignItems: 'center',
    alignSelf: 'stretch',
    flexDirection: 'row',
    gap: 8,
  },
  calenderIcon: {
    height: 16,
    width: 16,
  },
  cameraName: {
    color: colors.light.card.primaryText,
    fontSize: 16,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 24,
  },
  container: {
    alignItems: 'flex-start',
    alignSelf: 'stretch',
    backgroundColor: colors.light.card.bg,
    borderColor: colors.light.card.border,
    borderRadius: 14,
    borderStyle: 'solid',
    borderWidth: 1,
    gap: 12,
    padding: 16,
  },
  dateTime: {
    color: colors.light.other.dateTime,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 400,
    lineHeight: 20,
  },
  dateTimeContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: 8,
  },
  detailText: {
    color: colors.light.card.secondaryText,
    fontSize: 12,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 16,
  },
  downloadButton: {
    alignItems: 'center',
    backgroundColor: colors.light.button.bg,
    borderColor: colors.light.button.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
    height: 32,
    justifyContent: 'center',
    width: 32,
  },
  playButton: {
    alignItems: 'center',
    backgroundColor: colors.light.buttonBk.bg,
    borderRadius: 8,
    flex: 1,
    flexDirection: 'row',
    gap: 6,
    height: 32,
    justifyContent: 'center',
  },
  playButtonText: {
    color: colors.light.buttonBk.text,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
  primaryTextsContainer: {
    alignItems: 'center',
    alignSelf: 'stretch',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  secondaryTextsContainer: {
    alignItems: 'center',
    alignSelf: 'stretch',
    flexDirection: 'row',
    gap: 16,
    justifyContent: 'flex-start',
  },
});

export default RecCard;
