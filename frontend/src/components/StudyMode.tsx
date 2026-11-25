import type { Card } from '../types'

type StudyModeProps = {
  cards: Card[]
  studyIndex: number
  isAnswerVisible: boolean
  onClose: () => void
  onShowAnswer: () => void
  onPrevious: () => void
  onNext: () => void
}

function StudyMode({
  cards,
  studyIndex,
  isAnswerVisible,
  onClose,
  onShowAnswer,
  onPrevious,
  onNext,
}: StudyModeProps) {
  const currentCard = cards[studyIndex]

  return (
    <section className="study-mode">
      <button
        type="button"
        className="back-button"
        onClick={onClose}
      >
        Вернуться к карточкам
      </button>

      <p className="study-counter">
        Карточка {studyIndex + 1} из {cards.length}
      </p>

      <article className="study-card">
        <h3>{currentCard.english}</h3>

        {isAnswerVisible ? (
          <>
            <p className="thai-word">{currentCard.thai}</p>
            <p className="transcription">
              {currentCard.transcription}
            </p>
          </>
        ) : (
          <button
            type="button"
            className="flip-button"
            onClick={onShowAnswer}
          >
            Показать ответ
          </button>
        )}
      </article>

      <div className="study-controls">
        <button
          type="button"
          onClick={onPrevious}
          disabled={studyIndex === 0}
        >
          Предыдущая
        </button>

        <button
          type="button"
          onClick={onNext}
          disabled={studyIndex === cards.length - 1}
        >
          Следующая
        </button>
      </div>
    </section>
  )
}

export default StudyMode