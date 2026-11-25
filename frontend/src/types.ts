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