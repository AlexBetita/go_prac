import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";
import axios from "axios";

export interface User {
	id: string;
	email: string;
	provider: "local" | "google";
	created_at: number;
	updated_at: number;
}

interface AuthState {
	token: string | null;
	user: User | null;
	status: "idle" | "loading" | "failed";
}
const initialState: AuthState = {
	token: localStorage.getItem("jwt"),
	user: null,
	status: "idle",
};

export const loginThunk = createAsyncThunk<
	string,
	{ email: string; password: string },
	{ rejectValue: string }
>("auth/login", async (body, { rejectWithValue }) => {
	try {
		const { data } = await axios.post("/api/auth/login", body);
		return data.token;
	} catch (e: any) {
		return rejectWithValue(e.response?.data?.message ?? "Login failed");
	}
});

export const registerThunk = createAsyncThunk<
	string,
	{ email: string; password: string },
	{ rejectValue: string }
>("auth/register", async (body, { rejectWithValue }) => {
	try {
		const { data } = await axios.post("/api/auth/register", body);
		return data.token;
	} catch (e: any) {
		return rejectWithValue(e.response?.data?.message ?? "Register failed");
	}
});

export const fetchProfileThunk = createAsyncThunk<
	User,
	void,
	{ state: { auth: AuthState } }
>("auth/profile", async (_, { getState, rejectWithValue }) => {
	const jwt = getState().auth.token;
	if (!jwt) return rejectWithValue("No token");
	try {
		const { data } = await axios.get("/api/auth/profile", {
			headers: { Authorization: `Bearer ${jwt}` },
		});
		return data;
	} catch {
		return rejectWithValue("Unauthorized");
	}
});

const slice = createSlice({
	name: "auth",
	initialState,
	reducers: {
		logout(state) {
			state.token = null;
			state.user = null;
			localStorage.removeItem("jwt");
		},
		setToken(state, action: PayloadAction<string>) {
			state.token = action.payload;
			localStorage.setItem("jwt", action.payload);
		},
	},
	extraReducers: (builder) => {
		builder
			.addCase(loginThunk.fulfilled, (st, a) => {
				st.token = a.payload;
				localStorage.setItem("jwt", a.payload);
			})
			.addCase(registerThunk.fulfilled, (st, a) => {
				st.token = a.payload;
				localStorage.setItem("jwt", a.payload);
			})
			.addCase(fetchProfileThunk.fulfilled, (st, a) => {
				st.user = a.payload;
			})
			.addCase(fetchProfileThunk.rejected, (st) => {
				st.token = null;
				st.user = null;
				localStorage.removeItem("jwt");
			});
	},
});

export const { logout, setToken } = slice.actions;
export default slice.reducer;
