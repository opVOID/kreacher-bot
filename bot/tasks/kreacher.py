import os
import pickle
from bot import ins
from bot.helpers.handler import skip_current
from bot.helpers.queues import clear_queue, get_active_chats

dir = os.path.dirname(os.path.abspath(__file__))
queues = os.path.join(dir, "../dbs/queues.pkl")
actives = os.path.join(dir, "../dbs/actives.pkl")


@ins.on_audio_playout_ended
async def audio_ended(gc, source):
    print(f"audio ended: {source}")


@ins.on_video_playout_ended
async def video_ended(gc, source):
    print(f"video ended: {source}")


@ins.on_playout_ended
async def media_ended(gc, source, media_type):
    print(f"{media_type} ended: {source}")
