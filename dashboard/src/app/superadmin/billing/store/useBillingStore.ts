import { create } from 'zustand';

export interface PlanFeature {
  id: string;
  name: string;
}

export interface SubscriptionPlan {
  id: number;
  name: string;
  code: string;
  description: string;
  price: number;
  max_users: number;
  max_members: number;
  is_popular: boolean;
  // UI-only properties
  colorTheme?: 'blue' | 'purple' | 'amber';
}

interface BillingState {
  plans: SubscriptionPlan[];
  isLoading: boolean;
  error: string | null;
  fetchPlans: () => Promise<void>;
  addPlan: (plan: Omit<SubscriptionPlan, 'id'>) => Promise<void>;
  updatePlan: (id: number, plan: Partial<SubscriptionPlan>) => Promise<void>;
  deletePlan: (id: number) => Promise<void>;
}

const getAuthHeaders = () => {
  const token = localStorage.getItem("token");
  return {
    "Content-Type": "application/json",
    "Authorization": `Bearer ${token}`,
  };
};

export const useBillingStore = create<BillingState>((set, get) => ({
  plans: [],
  isLoading: false,
  error: null,

  fetchPlans: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch("/api/v1/billing/plans", {
        headers: getAuthHeaders(),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.message || "Failed to fetch plans");
      
      // Add UI-only color themes based on index, but use real is_popular
      const enrichedPlans = (result.data || []).map((plan: any, index: number) => ({
        ...plan,
        colorTheme: index === 0 ? 'blue' : index === 1 ? 'purple' : 'amber',
      }));
      
      set({ plans: enrichedPlans, isLoading: false });
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  addPlan: async (planData) => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch("/api/v1/billing/admin/plans", {
        method: "POST",
        headers: getAuthHeaders(),
        body: JSON.stringify(planData),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.message || "Failed to create plan");
      
      await get().fetchPlans();
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err;
    }
  },

  updatePlan: async (id, planData) => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch(`/api/v1/billing/admin/plans/${id}`, {
        method: "PUT",
        headers: getAuthHeaders(),
        body: JSON.stringify(planData),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.message || "Failed to update plan");
      
      await get().fetchPlans();
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err;
    }
  },

  deletePlan: async (id) => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch(`/api/v1/billing/admin/plans/${id}`, {
        method: "DELETE",
        headers: getAuthHeaders(),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.message || "Failed to delete plan");
      
      set((state) => ({
        plans: state.plans.filter((p) => p.id !== id),
        isLoading: false
      }));
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err;
    }
  },
}));
