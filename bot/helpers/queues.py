import os
from pkl import load_pkl, dump_pkl

dir = os.path.dirname(os.path.abspath(__file__))
queues = os.path.join(dir, "../dbs/queues.pkl")
actives = os.path.join(dir, "../dbs/actives.pkl")


async def get_active_chats() -> list:
    ACTIVE = load_pkl(actives, "rb", "list")
    return ACTIVE


def add_to_queue(chat, name, url, ref, type):
    QUEUE = load_pkl(queues, "rb", "dict")
    ACTIVE = load_pkl(actives, "rb", "list")
    try:
        if chat.id in QUEUE:
            QUEUE[chat.id].append([name, url, ref, type])
            dump_pkl(queues, "wb", QUEUE)
            return int(len(QUEUE[chat.id]) - 1)
        if chat.id not in ACTIVE:
            ACTIVE.append(chat.id)
            dump_pkl(actives, "wb", ACTIVE)
        QUEUE[chat.id] = [[name, url, ref, type]]
        dump_pkl(queues, "wb", QUEUE)
    except Exception as e:
        raise e


def get_queue(chat):
    QUEUE = load_pkl(queues, "rb", "dict")
    try:
        if chat.id in QUEUE:
            return QUEUE[chat.id]
        return 0
    except Exception as e:
        raise e


def pop_an_item(chat):
    QUEUE = load_pkl(queues, "rb", "dict")
    try:
        if chat.id not in QUEUE:
            return 0
        QUEUE[chat.id].pop(0)
        dump_pkl(queues, "wb", QUEUE)
        return 1
    except Exception as e:
        raise e


def clear_queue(chat):
    QUEUE = load_pkl(queues, "rb", "dict")
    ACTIVE = load_pkl(actives, "rb", "list")
    try:
        if chat.id not in QUEUE:
            return 0
        QUEUE.pop(chat.id)
        dump_pkl(queues, "wb", QUEUE)
        if chat.id in ACTIVE:
            ACTIVE.remove(chat.id)
            dump_pkl(actives, "wb", ACTIVE)
        return 1
    except Exception as e:
        raise e
