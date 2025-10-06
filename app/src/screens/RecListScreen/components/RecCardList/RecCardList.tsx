import RecCard from '@/screens/RecListScreen/components/RecCard/RecCard';
import { SearchContext } from '@/screens/RecListScreen/contexts/SearchContext';
import { useContext } from 'react';

const RecCardList = () => {
  const ctx = useContext(SearchContext);
  if (!ctx) {
    throw new Error('Search contextはProviderの中で使用してください。');
  }
  const { records } = ctx;

  return (
    <>
      {records.map(rec => (
        <RecCard record={rec} key={rec.url} />
      ))}
    </>
  );
};

export default RecCardList;
