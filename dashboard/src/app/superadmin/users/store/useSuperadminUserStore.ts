import { create } from 'zustand';

export interface Koperasi {
  id: number;
  name: string;
  email: string;
  phone: string;
  address: string;
  logo: string;
  plan: string;
  status: string;
  created_at: string;
  settings: Record<string, any>;
}

interface SuperadminUserState {
  koperasis: Koperasi[];
  searchQuery: string;
  isLoading: boolean;
  error: string | null;
  setSearchQuery: (query: string) => void;
  fetchKoperasis: () => Promise<void>;
  addKoperasi: (data: any) => Promise<void>;
  updateKoperasi: (id: number, data: any) => Promise<void>;
  deleteKoperasi: (id: number) => Promise<void>; // Currently unused in API, stubbed for UI compatibility
}

const getAuthHeaders = () => {
  const token = localStorage.getItem("token");
  return {
    "Content-Type": "application/json",
    "Authorization": `Bearer ${token}`,
  };
};

export const useSuperadminUserStore = create<SuperadminUserState>((set, get) => ({
  koperasis: [],
  searchQuery: '',
  isLoading: false,
  error: null,
  
  setSearchQuery: (query) => set({ searchQuery: query }),
  
  fetchKoperasis: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch("/api/v1/organizations", {
        headers: getAuthHeaders(),
      });
      const result = await response.json();
      
      if (!response.ok) {
        throw new Error(result.message || "Failed to fetch koperasis");
      }
      
      set({ koperasis: result.data || [], isLoading: false });
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  addKoperasi: async (data: any) => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch("/api/v1/organizations/onboard", {
        method: "POST",
        headers: { "Content-Type": "application/json" }, // Onboarding is public, no auth required by Bruno file
        body: JSON.stringify(data),
      });
      const result = await response.json();
      
      if (!response.ok) {
        throw new Error(result.message || "Failed to onboard koperasi");
      }

      await get().fetchKoperasis(); // Refresh the list
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err; // Re-throw to be handled by UI
    }
  },

  updateKoperasi: async (id: number, data: any) => {
    set({ isLoading: true, error: null });
    try {
      // Backend expects PATCH for settings/plan
      const response = await fetch(`/api/v1/organizations/${id}/settings`, {
        method: "PATCH",
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
      });
      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || "Failed to update koperasi");
      }

      await get().fetchKoperasis(); // Refresh
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err;
    }
  },

  deleteKoperasi: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      // Currently the backend does not have a delete organization endpoint in the Bruno collection.
      // We will just do a standard API call assuming one might exist, or simulate it.
      console.warn("Delete endpoint not explicitly defined. Implementing logic stub.");
      set((state) => ({ koperasis: state.koperasis.filter(k => k.id !== id), isLoading: false }));
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err;
    }
  }
}));
