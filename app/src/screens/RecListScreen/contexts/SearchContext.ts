import { DropdownItem } from '@/screens/RecListScreen/components/Dropdown/Dropdown';
import { Record } from '@/services/recorderApi';
import { createContext } from 'react';

type SearchContextType = {
  selectedCamera: string;
  setSelectedCamera: React.Dispatch<React.SetStateAction<string>>;
  cameras: DropdownItem[];
  setCameras: React.Dispatch<React.SetStateAction<DropdownItem[]>>;
  date: Date;
  setDate: React.Dispatch<React.SetStateAction<Date>>;
  records: Record[];
  setRecords: React.Dispatch<React.SetStateAction<Record[]>>;
};

export const SearchContext = createContext<SearchContextType | undefined>(undefined);
