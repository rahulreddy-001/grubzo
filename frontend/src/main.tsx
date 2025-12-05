import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { Provider } from "react-redux";
import { BrowserRouter } from "react-router";
import { NotificationProvider } from "./context/NotificationProvider.tsx";
import { AuthProvider } from "./context/AuthProvider.tsx";
import store from "./services/store.ts";

import "@radix-ui/themes/styles.css";
import { Theme } from "@radix-ui/themes";

createRoot(document.getElementById("root")!).render(
  <Theme accentColor="violet" radius="small">
    <BrowserRouter>
      <Provider store={store}>
        <AuthProvider>
          <NotificationProvider>
            <App />
          </NotificationProvider>
        </AuthProvider>
      </Provider>
    </BrowserRouter>
  </Theme>
);
