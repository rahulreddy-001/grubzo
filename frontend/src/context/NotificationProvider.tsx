import React, { createContext, useContext, useState, useCallback } from "react";
import type { ReactNode, ReactElement } from "react";

interface NotificationContextProps {
  showNotification: (
    Component: ReactElement<any>,
    props?: Record<string, any>
  ) => void;
  hideNotification: () => void;
}

const NotificationContext = createContext<NotificationContextProps | undefined>(
  undefined
);

export const useNotification = () => {
  const context = useContext(NotificationContext);
  if (!context)
    throw new Error("useNotification must be used within NotificationProvider");
  return context;
};

interface NotificationProviderProps {
  children: ReactNode;
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({
  children,
}) => {
  const [visible, setVisible] = useState(false);
  const [Component, setComponent] = useState<ReactElement<any> | null>(null);
  const [componentProps, setComponentProps] = useState<Record<string, any>>({});

  const showNotification = useCallback(
    (Comp: ReactElement<any>, props: Record<string, any> = {}) => {
      setComponent(Comp);
      setComponentProps(props);
      setVisible(true);
    },
    []
  );

  const hideNotification = useCallback(() => {
    setVisible(false);
    setComponent(null);
    setComponentProps({});
  }, []);

  const notificationWrapper = Component
    ? React.cloneElement(Component, {
        ...componentProps,
        visible,
        setVisible,
      })
    : null;

  return (
    <NotificationContext.Provider
      value={{ showNotification, hideNotification }}
    >
      {children}
      {notificationWrapper}
    </NotificationContext.Provider>
  );
};
