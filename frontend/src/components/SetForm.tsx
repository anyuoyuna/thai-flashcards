import type { FormEvent } from 'react'

type SetFormProps = {
  title: string
  isSaving: boolean
  error: string
  onTitleChange: (value: string) => void
  onSubmit: (event: FormEvent<HTMLFormElement>) => void
}

function SetForm({
  title,
  isSaving,
  error,
  onTitleChange,
  onSubmit,
}: SetFormProps) {
  return (
    <form className="set-form" onSubmit={onSubmit}>
      <h3>Создать набор</h3>

      <input
        value={title}
        onChange={(event) => onTitleChange(event.target.value)}
        placeholder="Например: Thai Basics"
        required
      />

      <button type="submit" disabled={isSaving}>
        {isSaving ? 'Сохранение...' : 'Создать набор'}
      </button>

      {error && <p className="error">{error}</p>}
    </form>
  )
}

export default SetForm