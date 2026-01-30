import { create } from 'zustand'
import type { Project, Collection, Endpoint, Environment } from '@/types'

interface ProjectState {
  currentProject: Project | null
  collections: Collection[]
  endpoints: Endpoint[]
  environments: Environment[]
  selectedCollection: string | null
  selectedEndpoint: string | null
  setCurrentProject: (project: Project | null) => void
  setCollections: (collections: Collection[]) => void
  setEndpoints: (endpoints: Endpoint[]) => void
  setEnvironments: (environments: Environment[]) => void
  setSelectedCollection: (id: string | null) => void
  setSelectedEndpoint: (id: string | null) => void
}

export const useProjectStore = create<ProjectState>((set) => ({
  currentProject: null,
  collections: [],
  endpoints: [],
  environments: [],
  selectedCollection: null,
  selectedEndpoint: null,
  setCurrentProject: (project) => set({ currentProject: project }),
  setCollections: (collections) => set({ collections }),
  setEndpoints: (endpoints) => set({ endpoints }),
  setEnvironments: (environments) => set({ environments }),
  setSelectedCollection: (id) => set({ selectedCollection: id }),
  setSelectedEndpoint: (id) => set({ selectedEndpoint: id }),
}))
