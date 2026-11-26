import { useState, type FormEvent } from 'react'

type CsvImportFormProps = {
  isImporting: boolean
  error: string
  onImport: (file: File) => void
}

function CsvImportForm({
  isImporting,
  error,
  onImport,
}: CsvImportFormProps) {
  const [file, setFile] = useState<File | null>(null)

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (file === null) {
      return
    }

    onImport(file)
  }

  return (
    <form className="csv-form" onSubmit={handleSubmit}>
      <h3>Импорт карточек из CSV</h3>

      <input
        type="file"
        accept=".csv,text/csv"
        onChange={(event) => {
          setFile(event.target.files?.[0] ?? null)
        }}
      />

      <button
        type="submit"
        disabled={file === null || isImporting}
      >
        {isImporting ? 'Импорт...' : 'Импортировать'}
      </button>

      {error && <p className="error">{error}</p>}
    </form>
  )
}

export default CsvImportForm