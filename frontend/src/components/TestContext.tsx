import React, { createContext, useState, useContext, ReactNode } from 'react';

type TestContextType = {
  isTestActive: boolean;
  setIsTestActive: (active: boolean) => void;
};

const TestContext = createContext<TestContextType | undefined>(undefined);

export const TestProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isTestActive, setIsTestActive] = useState(false);

  return (
    <TestContext.Provider value={{ isTestActive, setIsTestActive }}>
      {children}
    </TestContext.Provider>
  );
};

export const useTest = () => {
  const context = useContext(TestContext);
  if (context === undefined) {
    throw new Error('useTest must be used within a TestProvider');
  }
  return context;
};
