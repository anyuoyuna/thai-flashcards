type CardEditFormProps = {
  english: string
  thai: string
  transcription: string
  error: string
  onEnglishChange: (value: string) => void
  onThaiChange: (value: string) => void
  onTranscriptionChange: (value: string) => void
  onSave: () => void
  onCancel: () => void
}

function CardEditForm({
  english,
  thai,
  transcription,
  error,
  onEnglishChange,
  onThaiChange,
  onTranscriptionChange,
  onSave,
  onCancel,
}: CardEditFormProps) {
  return (
    <div className="edit-card-form">
      <input
        value={english}
        onChange={(event) => onEnglishChange(event.target.value)}
      />

      <input
        value={thai}
        onChange={(event) => onThaiChange(event.target.value)}
      />

      <input
        value={transcription}
        onChange={(event) => onTranscriptionChange(event.target.value)}
      />

      <button type="button" onClick={onSave}>
        Сохранить
      </button>

      <button type="button" onClick={onCancel}>
        Отмена
      </button>

      {error && <p className="error">{error}</p>}
    </div>
  )
}

export default CardEditForm