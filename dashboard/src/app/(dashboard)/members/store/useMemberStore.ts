import { create } from 'zustand';

export type MemberStatus = 'active' | 'inactive';

export interface Member {
  id: string;
  name: string;
  email: string;
  phone: string;
  status: MemberStatus;
  joinDate: string;
  balance: number;
}

interface MemberState {
  members: Member[];
  searchQuery: string;
  statusFilter: 'all' | MemberStatus;
  setSearchQuery: (query: string) => void;
  setStatusFilter: (status: 'all' | MemberStatus) => void;
  addMember: (member: Omit<Member, 'id' | 'joinDate'>) => void;
  updateMember: (id: string, member: Partial<Member>) => void;
  deleteMember: (id: string) => void;
}

const mockMembers: Member[] = [
  { id: '1', name: 'Budi Santoso', email: 'budi@example.com', phone: '081234567890', status: 'active', joinDate: '2023-01-15', balance: 5000000 },
  { id: '2', name: 'Siti Aminah', email: 'siti@example.com', phone: '081298765432', status: 'active', joinDate: '2023-03-22', balance: 12500000 },
  { id: '3', name: 'Andi Wijaya', email: 'andi@example.com', phone: '085612349876', status: 'inactive', joinDate: '2022-11-10', balance: 0 },
  { id: '4', name: 'Lestari Ningsih', email: 'lestari@example.com', phone: '081345678901', status: 'active', joinDate: '2023-05-05', balance: 7500000 },
  { id: '5', name: 'Rahmat Hidayat', email: 'rahmat@example.com', phone: '082156781234', status: 'active', joinDate: '2023-08-12', balance: 2100000 },
  { id: '6', name: 'Dewi Sartika', email: 'dewi@example.com', phone: '085789012345', status: 'inactive', joinDate: '2022-12-01', balance: 150000 },
  { id: '7', name: 'Hendro Siswanto', email: 'hendro@example.com', phone: '081245670987', status: 'active', joinDate: '2024-01-20', balance: 9000000 },
  { id: '8', name: 'Maya Sari', email: 'maya@example.com', phone: '081398761234', status: 'active', joinDate: '2023-10-10', balance: 4300000 },
  { id: '9', name: 'Eko Prasetyo', email: 'eko@example.com', phone: '081122334455', status: 'active', joinDate: '2023-11-05', balance: 3200000 },
  { id: '10', name: 'Sri Wahyuni', email: 'sri@example.com', phone: '081566778899', status: 'active', joinDate: '2024-02-14', balance: 11000000 },
  { id: '11', name: 'Agus Gunawan', email: 'agus@example.com', phone: '081900112233', status: 'inactive', joinDate: '2023-07-20', balance: 50000 },
  { id: '12', name: 'Nina Wati', email: 'nina@example.com', phone: '082233445566', status: 'active', joinDate: '2024-03-01', balance: 8400000 },
];

export const useMemberStore = create<MemberState>((set) => ({
  members: mockMembers,
  searchQuery: '',
  statusFilter: 'all',
  setSearchQuery: (query) => set({ searchQuery: query }),
  setStatusFilter: (status) => set({ statusFilter: status }),
  addMember: (memberData) => set((state) => ({
    members: [
      {
        ...memberData,
        id: Math.random().toString(36).substr(2, 9),
        joinDate: new Date().toISOString().split('T')[0],
      },
      ...state.members,
    ]
  })),
  updateMember: (id, memberData) => set((state) => ({
    members: state.members.map((m) => m.id === id ? { ...m, ...memberData } : m)
  })),
  deleteMember: (id) => set((state) => ({
    members: state.members.filter((m) => m.id !== id)
  }))
}));
