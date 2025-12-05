import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./auth/auth.slice";
import itemReducer from "./item/item.slice";
import rbacReducer from "./common/rbac.slice";
import commonReducer from "./common/common.slice";
import employeeReducer from "./common/employee.slice";

const store = configureStore({
  reducer: {
    auth: authReducer,
    item: itemReducer,
    rbac: rbacReducer,
    common: commonReducer,
    emp: employeeReducer,
  },
});

export default store;
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
