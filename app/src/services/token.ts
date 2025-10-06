import storage from './storage';

const TOKEN_STORAGE_KEY = 'token';

const loadToken = async (): Promise<string> => {
  try {
    const token: string = await storage.load({
      key: TOKEN_STORAGE_KEY,
      autoSync: false,
    });
    return token;
  } catch (e: unknown) {
    if (e instanceof Error) {
      if (e.name === 'NotFoundError' || e.name === 'ExpiredError') {
        return '';
      }
      throw e;
    }
  }
  return '';
};

const saveToken = async (token: string): Promise<void> => {
  try {
    await storage.save({
      key: TOKEN_STORAGE_KEY,
      data: token,
      expires: null,
    });
  } catch (e: unknown) {
    if (e instanceof Error) {
      throw e;
    }
  }
};

export { loadToken, saveToken };
