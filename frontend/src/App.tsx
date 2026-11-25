import { useEffect, useState, type FormEvent } from 'react'
import './App.css'
import StudyMode from './components/StudyMode'
import type { Card, CardSet } from './types'
import CardItem from './components/CardItem'
import CardEditForm from './components/CardEditForm'
import CardForm from './components/CardForm'
import SetForm from './components/SetForm'
import {
  createCard as createCardRequest,
  createSet as createSetRequest,
  deleteCard as deleteCardRequest,
  getCards,
  getSets,
  updateCard as updateCardRequest,
} from './api'

function App() {
  const [cardSets, setCardSets] = useState<CardSet[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [selectedSetId, setSelectedSetId] = useState<number | null>(null)
  const [cards, setCards] = useState<Card[]>([])
  const [isCardsLoading, setIsCardsLoading] = useState(false)
  const [cardsError, setCardsError] = useState('')

  const [english, setEnglish] = useState('')
  const [thai, setThai] = useState('')
  const [transcription, setTranscription] = useState('')
  const [isSavingCard, setIsSavingCard] = useState(false)
  const [cardFormError, setCardFormError] = useState('')

  const [setTitle, setSetTitle] = useState('')
  const [isSavingSet, setIsSavingSet] = useState(false)
  const [setFormError, setSetFormError] = useState('')

  const [editingCardId, setEditingCardId] = useState<number | null>(null)
  const [editEnglish, setEditEnglish] = useState('')
  const [editThai, setEditThai] = useState('')
  const [editTranscription, setEditTranscription] = useState('')
  const [editFormError, setEditFormError] = useState('')

  const [isStudyMode, setIsStudyMode] = useState(false)
  const [studyIndex, setStudyIndex] = useState(0)
  const [isAnswerVisible, setIsAnswerVisible] = useState(false)

  useEffect(() => {
    async function loadCardSets() {
      try {
        const data = await getSets()
        setCardSets(data)
      } catch {
        setError('Не удалось подключиться к серверу')
      } finally {
        setIsLoading(false)
      }
    }

    loadCardSets()
  }, [])

  async function createSet(event: FormEvent<HTMLFormElement>) {
      event.preventDefault()

      setIsSavingSet(true)
      setSetFormError('')

      try {
        const createdSet = await createSetRequest(setTitle)

        setCardSets((currentSets) => [...currentSets, createdSet])
        setSetTitle('')
      } catch {
        setSetFormError('Не удалось создать набор')
      } finally {
        setIsSavingSet(false)
      }
  }

  async function openSet(setId: number) {
    setSelectedSetId(setId)
    setIsCardsLoading(true)
    setCardsError('')

    try {
      const data = await getCards(setId)
      setCards(data)
    } catch {
      setCardsError('Не удалось загрузить карточки')
      setCards([])
    } finally {
      setIsCardsLoading(false)
    }
  }

  function startEditingCard(card: Card) {
    setEditingCardId(card.id)
    setEditEnglish(card.english)
    setEditThai(card.thai)
    setEditTranscription(card.transcription)
    setEditFormError('')
  }

  function cancelEditingCard() {
    setEditingCardId(null)
    setEditFormError('')
  }

  async function updateCard(cardId: number) {
    setEditFormError('')

    try {
      const updatedCard = await updateCardRequest(cardId, {
        english: editEnglish,
        thai: editThai,
        transcription: editTranscription,
      })

      setCards((currentCards) =>
        currentCards.map((card) =>
          card.id === cardId ? updatedCard : card,
        ),
      )

      setEditingCardId(null)
    } catch {
      setEditFormError('Не удалось обновить карточку')
    }
  }

  async function createCard(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (selectedSetId === null) {
      return
    }

    setIsSavingCard(true)
    setCardFormError('')

    try {
      const createdCard = await createCardRequest(selectedSetId, {
        english,
        thai,
        transcription,
      })

      setCards((currentCards) => [...currentCards, createdCard])

      setEnglish('')
      setThai('')
      setTranscription('')
    } catch {
      setCardFormError('Не удалось создать карточку')
    } finally {
      setIsSavingCard(false)
    }
  }

  async function deleteCard(cardId: number) {
    const shouldDelete = window.confirm('Удалить эту карточку?')

    if (!shouldDelete) {
      return
    }

    try {
      await deleteCardRequest(cardId)

      setCards((currentCards) =>
        currentCards.filter((card) => card.id !== cardId),
      )
    } catch {
      setCardsError('Не удалось удалить карточку')
    }
  }

  function closeSet() {
    setSelectedSetId(null)
    setCards([])
    setCardsError('')
    setIsStudyMode(false)
    setStudyIndex(0)
    setIsAnswerVisible(false)
  }

  function startStudyMode() {
    setStudyIndex(0)
    setIsAnswerVisible(false)
    setIsStudyMode(true)
  }

  function closeStudyMode() {
    setIsStudyMode(false)
    setIsAnswerVisible(false)
  }

  function showNextCard() {
    if (studyIndex < cards.length - 1) {
      setStudyIndex((currentIndex) => currentIndex + 1)
      setIsAnswerVisible(false)
    }
  }

  function showPreviousCard() {
    if (studyIndex > 0) {
      setStudyIndex((currentIndex) => currentIndex - 1)
      setIsAnswerVisible(false)
    }
  }

  const selectedSet = cardSets.find((set) => set.id === selectedSetId)

  return (
    <main className="app">
      <header className="app-header">
        <h1>Thai Flashcards</h1>
        <p>Наборы карточек для изучения тайского языка</p>
      </header>

      {selectedSetId === null ? (
        <section>
          <h2>Мои наборы</h2>

          <SetForm
            title={setTitle}
            isSaving={isSavingSet}
            error={setFormError}
            onTitleChange={setSetTitle}
            onSubmit={createSet}
          />

          {isLoading && <p>Загрузка наборов...</p>}

          {error && <p className="error">{error}</p>}

          {!isLoading && !error && cardSets.length === 0 && (
            <p>Наборов пока нет</p>
          )}

          <div className="sets-list">
            {cardSets.map((set) => (
              <article className="set-card" key={set.id}>
                <h3>{set.title}</h3>
                <button type="button" onClick={() => openSet(set.id)}>
                  Открыть
                </button>
              </article>
            ))}
          </div>
        </section>
      ) : (
        <section>
          <button type="button" className="back-button" onClick={closeSet}>
            ← Назад к наборам
          </button>

          <h2>{selectedSet?.title}</h2>

          {cards.length > 0 && !isStudyMode && (
            <button
              type="button"
              className="study-button"
              onClick={startStudyMode}
            >
              Начать изучение
            </button>
          )}

          {!isStudyMode && (
            <CardForm
              english={english}
              thai={thai}
              transcription={transcription}
              isSaving={isSavingCard}
              error={cardFormError}
              onEnglishChange={setEnglish}
              onThaiChange={setThai}
              onTranscriptionChange={setTranscription}
              onSubmit={createCard}
            />
          )}
          {isCardsLoading && <p>Загрузка карточек...</p>}

          {cardsError && <p className="error">{cardsError}</p>}

          {!isCardsLoading && !cardsError && cards.length === 0 && (
            <p>В этом наборе пока нет карточек</p>
          )}
          {isStudyMode ? (
            <StudyMode
              cards={cards}
              studyIndex={studyIndex}
              isAnswerVisible={isAnswerVisible}
              onClose={closeStudyMode}
              onShowAnswer={() => setIsAnswerVisible(true)}
              onPrevious={showPreviousCard}
              onNext={showNextCard}
            />
          ) : (
            <div className="cards-list">
            {cards.map((card) => (
              <article className="card" key={card.id}>
                {editingCardId === card.id ? (
                  <CardEditForm
                    english={editEnglish}
                    thai={editThai}
                    transcription={editTranscription}
                    error={editFormError}
                    onEnglishChange={setEditEnglish}
                    onThaiChange={setEditThai}
                    onTranscriptionChange={setEditTranscription}
                    onSave={() => updateCard(card.id)}
                    onCancel={cancelEditingCard}
                  />
                ) : (
                  <CardItem
                    card={card}
                    onEdit={startEditingCard}
                    onDelete={deleteCard}
                  />
                )}
              </article>
            ))}
          </div>
          )}
        </section>
      )}
    </main>
  )
}

export default App