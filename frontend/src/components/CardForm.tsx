import type { FormEvent } from 'react'

type CardFormProps = {
  english: string
  thai: string
  transcription: string
  isSaving: boolean
  error: string
  onEnglishChange: (value: string) => void
  onThaiChange: (value: string) => void
  onTranscriptionChange: (value: string) => void
  onSubmit: (event: FormEvent<HTMLFormElement>) => void
}

function CardForm({
  english,
  thai,
  transcription,
  isSaving,
  error,
  onEnglishChange,
  onThaiChange,
  onTranscriptionChange,
  onSubmit,
}: CardFormProps) {
  return (
    <form className="card-form" onSubmit={onSubmit}>
      <h3>Добавить карточку</h3>

      <label>
        English
        <input
          value={english}
          onChange={(event) => onEnglishChange(event.target.value)}
          placeholder="hello"
          required
        />
      </label>

      <label>
        Thai
        <input
          value={thai}
          onChange={(event) => onThaiChange(event.target.value)}
          placeholder="สวัสดี"
          required
        />
      </label>

      <label>
        Транскрипция
        <input
          value={transcription}
          onChange={(event) => onTranscriptionChange(event.target.value)}
          placeholder="sawasdee"
          required
        />
      </label>

      <button type="submit" disabled={isSaving}>
        {isSaving ? 'Сохранение...' : 'Добавить карточку'}
      </button>

      {error && <p className="error">{error}</p>}
    </form>
  )
}

export default CardForm