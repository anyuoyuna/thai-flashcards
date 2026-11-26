export type CardSet = {
  id: number
  title: string
}

export type Card = {
  id: number
  setId: number
  english: string
  thai: string
  transcription: string
}

export type CardProgress = {
  cardId: number
  correctCount: number
  incorrectCount: number
  lastReviewedAt: string
}

export type SetProgress = {
  setId: number
  totalCards: number
  studiedCards: number
  correctCount: number
  incorrectCount: number
}