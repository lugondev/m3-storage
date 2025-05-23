import { create } from 'zustand'

interface BreadcrumbLabelState {
	labels: Record<string, string>
	setLabel: (id: string, label: string) => void
	getLabel: (id: string) => string | undefined
}

export const useBreadcrumbLabelStore = create<BreadcrumbLabelState>((set, get) => ({
	labels: {},
	setLabel: (id, label) =>
		set((state) => ({
			labels: {
				...state.labels,
				[id]: label,
			},
		})),
	getLabel: (id) => get().labels[id],
}))
