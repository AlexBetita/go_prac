import {
  persistStore,
  persistReducer,
  FLUSH,
  PAUSE,
  PERSIST,
  PURGE,
  REGISTER,
  REHYDRATE,
} from "redux-persist";
import storage from "redux-persist/lib/storage";

import { configureStore } from "@reduxjs/toolkit";

import authReducer from "./slices/authSlice";
import postReducer from "./slices/postSlice";
import botReducer from "./slices/botSlice";

const authPersistConfig = {
  key: "auth",
  storage,
  whitelist: ["token", "tokenExp", "user"],
};

const postPersistConfig = {
  key: "posts",
  storage,
  whitelist: ["posts"]
}

const persistedAuthReducer = persistReducer(authPersistConfig, authReducer);
const persistedPostReducer = persistReducer(postPersistConfig, postReducer);

export const store = configureStore({
  reducer: {
    auth: persistedAuthReducer,
    posts: persistedPostReducer,
    bot: botReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
      },
    }),
});

export const persistor = persistStore(store);
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;