import type { Card } from '../types'

type CardItemProps = {
  card: Card
  onEdit: (card: Card) => void
  onDelete: (cardId: number) => void
}

function CardItem({ card, onEdit, onDelete }: CardItemProps) {
  return (
    <>
      <h3>{card.english}</h3>

      <p className="thai-word">{card.thai}</p>

      <p className="transcription">
        {card.transcription}
      </p>

      <button
        type="button"
        className="edit-button"
        onClick={() => onEdit(card)}
      >
        Редактировать
      </button>

      <button
        type="button"
        className="delete-button"
        onClick={() => onDelete(card.id)}
      >
        Удалить
      </button>
    </>
  )
}

export default CardItem