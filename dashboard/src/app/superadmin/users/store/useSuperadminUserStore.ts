import { create } from 'zustand';

export type KoperasiStatus = 'active' | 'inactive' | 'trial';
export type KoperasiPlan = 'Enterprise' | 'Business' | 'Starter';

export interface Koperasi {
  id: string;
  name: string;
  location: string;
  email: string;
  plan: KoperasiPlan;
  status: KoperasiStatus;
  joinDate: string;
}

interface SuperadminUserState {
  koperasis: Koperasi[];
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  addKoperasi: (koperasi: Omit<Koperasi, 'id' | 'joinDate'>) => void;
  updateKoperasi: (id: string, koperasi: Partial<Koperasi>) => void;
  deleteKoperasi: (id: string) => void;
}

const mockKoperasis: Koperasi[] = [
  { id: '1', name: 'Koperasi Maju Jaya', location: 'Jakarta, Indonesia', email: 'admin@majujaya.com', plan: 'Enterprise', status: 'active', joinDate: 'Jan 12, 2024' },
  { id: '2', name: 'Kud Binangun', location: 'Yogyakarta, Indonesia', email: 'contact@binangun.org', plan: 'Business', status: 'active', joinDate: 'Feb 05, 2024' },
  { id: '3', name: 'Kopka Mandiri', location: 'Bandung, Indonesia', email: 'info@kopkamandiri.com', plan: 'Starter', status: 'trial', joinDate: 'Mar 15, 2024' },
];

export const useSuperadminUserStore = create<SuperadminUserState>((set) => ({
  koperasis: mockKoperasis,
  searchQuery: '',
  setSearchQuery: (query) => set({ searchQuery: query }),
  addKoperasi: (koperasiData) => set((state) => ({
    koperasis: [
      {
        ...koperasiData,
        id: Math.random().toString(36).substr(2, 9),
        joinDate: new Date().toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' }),
      },
      ...state.koperasis,
    ]
  })),
  updateKoperasi: (id, koperasiData) => set((state) => ({
    koperasis: state.koperasis.map((k) => k.id === id ? { ...k, ...koperasiData } : k)
  })),
  deleteKoperasi: (id) => set((state) => ({
    koperasis: state.koperasis.filter((k) => k.id !== id)
  }))
}));
