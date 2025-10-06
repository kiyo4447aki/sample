/* eslint-disable @typescript-eslint/no-require-imports */

const { withAndroidManifest } = require('@expo/config-plugins');

function ensureToolsNamespace(manifest) {
  manifest.manifest.$ = manifest.manifest.$ || {};
  if (!manifest.manifest.$['xmlns:tools']) {
    manifest.manifest.$['xmlns:tools'] = 'http://schemas.android.com/tools';
  }
}

function upsertDefaultChannelMeta(androidManifest, channelId = 'alerts') {
  const app = androidManifest.manifest.application?.[0];
  if (!app) return androidManifest;

  app['meta-data'] = app['meta-data'] || [];

  const NAME = 'com.google.firebase.messaging.default_notification_channel_id';

  const idx = app['meta-data'].findIndex(m => m.$ && m.$['android:name'] === NAME);

  const entry = {
    $: {
      'android:name': NAME,
      'android:value': channelId,
      'tools:replace': 'android:value', // ← これが鍵
    },
  };

  if (idx >= 0) {
    app['meta-data'][idx] = { $: { ...app['meta-data'][idx].$, ...entry.$ } };
  } else {
    app['meta-data'].push(entry);
  }

  return androidManifest;
}

const withFcmDefaultChannel = (config, { channelId = 'alerts' } = {}) =>
  withAndroidManifest(config, c => {
    ensureToolsNamespace(c.modResults);
    c.modResults = upsertDefaultChannelMeta(c.modResults, channelId);
    return c;
  });

module.exports = withFcmDefaultChannel;
