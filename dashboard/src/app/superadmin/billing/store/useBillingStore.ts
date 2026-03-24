import { create } from 'zustand';

export interface PlanFeature {
  id: string;
  name: string;
}

export interface SubscriptionPlan {
  id: string;
  name: string;
  priceAmount: number;
  priceDisplay: string;
  interval: string;
  description: string;
  features: PlanFeature[];
  isPopular: boolean;
  colorTheme: 'blue' | 'purple' | 'amber';
}

interface BillingState {
  plans: SubscriptionPlan[];
  addPlan: (plan: Omit<SubscriptionPlan, 'id'>) => void;
  updatePlan: (id: string, plan: Partial<SubscriptionPlan>) => void;
  deletePlan: (id: string) => void;
}

const mockPlans: SubscriptionPlan[] = [
  {
    id: '1',
    name: 'Starter',
    priceAmount: 499000,
    priceDisplay: 'Rp 499k',
    interval: '/mo',
    description: 'Up to 100 members & basic accounting',
    features: [
      { id: 'f1', name: 'Core Banking Module' },
      { id: 'f2', name: 'Member Management' },
      { id: 'f3', name: 'Simple Reporting' }
    ],
    isPopular: false,
    colorTheme: 'blue'
  },
  {
    id: '2',
    name: 'Business',
    priceAmount: 1490000,
    priceDisplay: 'Rp 1.49M',
    interval: '/mo',
    description: 'Up to 1,000 members & advanced features',
    features: [
      { id: 'f4', name: 'Everything in Starter' },
      { id: 'f5', name: 'Inventory Module' },
      { id: 'f6', name: 'Advanced Financial Reports' }
    ],
    isPopular: true,
    colorTheme: 'purple'
  },
  {
    id: '3',
    name: 'Enterprise',
    priceAmount: 0,
    priceDisplay: 'Custom',
    interval: '',
    description: 'Unlimited members & custom modules',
    features: [
      { id: 'f7', name: 'Full Suite Access' },
      { id: 'f8', name: 'Multi-branch Support' },
      { id: 'f9', name: 'Priority 24/7 Support' }
    ],
    isPopular: false,
    colorTheme: 'amber'
  }
];

export const useBillingStore = create<BillingState>((set) => ({
  plans: mockPlans,
  addPlan: (planData) => set((state) => ({
    plans: [
      ...state.plans,
      {
        ...planData,
        id: Math.random().toString(36).substr(2, 9),
      }
    ]
  })),
  updatePlan: (id, planData) => set((state) => ({
    plans: state.plans.map((p) => p.id === id ? { ...p, ...planData } : p)
  })),
  deletePlan: (id) => set((state) => ({
    plans: state.plans.filter((p) => p.id !== id)
  }))
}));
