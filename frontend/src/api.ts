import type { 
  Card, 
  CardSet, 
  CardProgress,
  SetProgress,
} from './types'

export async function getSets(): Promise<CardSet[]> {
  const response = await fetch('/api/sets')

  if (!response.ok) {
    throw new Error('Failed to load sets')
  }

  return response.json()
}

export async function createSet(title: string): Promise<CardSet> {
  const response = await fetch('/api/sets', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ title }),
  })

  if (!response.ok) {
    throw new Error('Failed to create set')
  }

  return response.json()
}

export async function getCards(setId: number): Promise<Card[]> {
  const response = await fetch(`/api/sets/${setId}/cards`)

  if (!response.ok) {
    throw new Error('Failed to load cards')
  }

  const data: Card[] | null = await response.json()

  return data ?? []
}

export async function createCard(
  setId: number,
  cardData: {
    english: string
    thai: string
    transcription: string
  },
): Promise<Card> {
  const response = await fetch(`/api/sets/${setId}/cards`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(cardData),
  })

  if (!response.ok) {
    throw new Error('Failed to create card')
  }

  return response.json()
}

export async function updateCard(
  cardId: number,
  cardData: {
    english: string
    thai: string
    transcription: string
  },
): Promise<Card> {
  const response = await fetch(`/api/cards/${cardId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(cardData),
  })

  if (!response.ok) {
    throw new Error('Failed to update card')
  }

  return response.json()
}

export async function deleteCard(cardId: number): Promise<void> {
  const response = await fetch(`/api/cards/${cardId}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    throw new Error('Failed to delete card')
  }
}

export async function reviewCard(
  cardId: number,
  known: boolean,
): Promise<CardProgress> {
  const response = await fetch(`/api/cards/${cardId}/review`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ known }),
  })

  if (!response.ok) {
    throw new Error('Failed to save card review')
  }

  return response.json()
}

export async function getCardProgress(
  cardId: number,
): Promise<CardProgress> {
  const response = await fetch(`/api/cards/${cardId}/progress`)

  if (!response.ok) {
    throw new Error('Failed to load card progress')
  }

  return response.json()
}

export async function getSetProgress(
  setId: number,
): Promise<SetProgress> {
  const response = await fetch(`/api/sets/${setId}/progress`)

  if (!response.ok) {
    throw new Error('Failed to load set progress')
  }

  return response.json()
}

export async function importCards(
  setId: number,
  file: File,
): Promise<Card[]> {
  const formData = new FormData()
  formData.append('file', file)

  const response = await fetch(`/api/sets/${setId}/cards/import`, {
    method: 'POST',
    body: formData,
  })

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || 'Failed to import cards')
  }

  return response.json()
}