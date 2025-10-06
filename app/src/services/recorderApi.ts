import axios from 'axios';

type Record = {
  name: string;
  url: string;
  timestamp: string;
  size: number;
};

type RecordListResponse = {
  status: string;
  records: Record[];
  const: number;
};

const initclient = (recorderId: string, token: string) => {
  const url = ``;
  const client = axios.create({
    baseURL: url,
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return client;
};

export { Record, RecordListResponse };

export default initclient;
