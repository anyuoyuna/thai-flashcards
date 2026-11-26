import type { Card, CardProgress} from '../types'

type StudyModeProps = {
  cards: Card[]
  studyIndex: number
  isAnswerVisible: boolean
  onClose: () => void
  onShowAnswer: () => void
  onPrevious: () => void
  onNext: () => void
  onReview: (known: boolean) => void
  progress: CardProgress | null
  isReviewing: boolean
}

function StudyMode({
  cards,
  studyIndex,
  isAnswerVisible,
  progress,
  onClose,
  onShowAnswer,
  onPrevious,
  onNext,
  onReview,
  isReviewing,
}: StudyModeProps) {
  const currentCard = cards[studyIndex]

  return (
    <section className="study-mode">
      <button
        type="button"
        className="back-button"
        onClick={onClose}
        disabled={isReviewing}
      >
        Вернуться к карточкам
      </button>

      <p className="study-counter">
        Карточка {studyIndex + 1} из {cards.length}
      </p>
      {progress && (
        <p className="progress-info">
          Знаю: {progress.correctCount} · Повторить:{' '}
          {progress.incorrectCount}
        </p>
      )}
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
        {isAnswerVisible && (
          <div className="review-controls">
            <button
              type="button"
              onClick={() => onReview(false)}
              disabled={isReviewing}
            >
              Повторить
            </button>

            <button
              type="button"
              onClick={() => onReview(true)}
              disabled={isReviewing}
            >
              Знаю
            </button>
          </div>
        )}
      </article>

      <div className="study-controls">
        <button
          type="button"
          onClick={onPrevious}
          disabled={studyIndex === 0 || isReviewing}
        >
          Предыдущая
        </button>

        <button
          type="button"
          onClick={onNext}
          disabled={studyIndex === cards.length - 1 || isReviewing}
        >
          Следующая
        </button>
      </div>
    </section>
  )
}

export default StudyMode