import store from "../store";
import { loginUser, signupUser, fetchUser, logoutUser } from "./auth.slice";
import type {
  SignupRequest,
  LoginRequest,
  User,
  LoginResponse,
  SignupResponse,
} from "../../types/auth";

const AuthService = {
  async isAuthenticated(): Promise<boolean> {
    try {
      const state = store.getState();
      const user = state.auth?.user || null;
      return !!user;
    } catch (err) {
      console.error("Error checking auth from store", err);
      return false;
    }
  },

  async login(data: LoginRequest): Promise<LoginResponse> {
    return store.dispatch(loginUser(data)).unwrap();
  },

  async signup(userData: SignupRequest): Promise<SignupResponse> {
    return store.dispatch(signupUser(userData)).unwrap();
  },

  async fetchUser(): Promise<User> {
    return store.dispatch(fetchUser()).unwrap();
  },

  async logout(): Promise<void> {
    await store.dispatch(logoutUser()).unwrap();
  },
};

export default AuthService;
