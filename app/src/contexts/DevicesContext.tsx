import { useAuth } from '@/contexts/AuthContext';
import api from '@/services/api';
import React, {
  createContext,
  PropsWithChildren,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';

type ConnectionInfoResp = {
  device_id: string;
  janus_password: string;
  janus_ws: string;
  status: string;
};

type ConnectionInfoType = {
  deviceId: string;
  password: string;
  janusUrl: string;
};

type DeviceListResp = {
  devices: string[];
};

type DevicesContextType = {
  deviceList: string[];
  device: string;
  setDevice: React.Dispatch<React.SetStateAction<string>>;
  fetchDeviceList: () => void;
  getConnectionInfo: () => Promise<ConnectionInfoType | undefined>;
};

const DevicesContext = createContext<DevicesContextType | undefined>(undefined);

const useDevices = () => {
  const ctx = useContext(DevicesContext);
  if (!ctx) throw new Error('useDevicesはAuthProviderの中で使用してください');
  return ctx;
};

const DevicesProvider = ({ children }: PropsWithChildren) => {
  const [device, setDevice] = useState<string>('');
  const [deviceList, setDeviceList] = useState<string[]>([]);

  const { isLoggedIn } = useAuth();

  //TODO 複数デバイス対応
  useEffect(() => {
    if (deviceList && deviceList.length > 0) {
      setDevice(deviceList[0]);
    }
  }, [deviceList]);

  const fetchDeviceList = async () => {
    if (!isLoggedIn) return;
    const res = await api.get<DeviceListResp>('/devices');
    const devices = res.data;
    if (devices.devices) setDeviceList(devices.devices);
  };

  const getConnectionInfo = async () => {
    if (!isLoggedIn) return;
    const res = await api.get<ConnectionInfoResp>(`/devices/${device}/connection-info`);
    const connectionInfo: ConnectionInfoType | undefined = {
      deviceId: res.data.device_id,
      password: res.data.janus_password,
      janusUrl: res.data.janus_ws,
    };
    return connectionInfo;
  };

  const value = useMemo(
    () => ({
      deviceList,
      device,
      setDevice,
      fetchDeviceList,
      getConnectionInfo,
    }),
    [deviceList, device, setDevice, fetchDeviceList, getConnectionInfo],
  );

  return <DevicesContext.Provider value={value}>{children}</DevicesContext.Provider>;
};

export { ConnectionInfoType, DevicesProvider, useDevices };
