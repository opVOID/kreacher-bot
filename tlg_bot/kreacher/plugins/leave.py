from kreacher import kreacher
from kreacher.helpers.voice_chats import VOICE_CHATS
from telethon import events


@kreacher.on(events.NewMessage(pattern="[!?/]leave"))
async def leave_handler(event):
    try:
        chat = await event.get_chat()
        proto = VOICE_CHATS[chat.id]
        if proto is None:
            raise Exception("Streaming is not active")
        await proto.leave_current_group_call()
        VOICE_CHATS.pop(chat.id)
        await event.reply("<i>Goodbye master, just call me if you need me. \n\nVoice Chat left successfully.</i>", parse_mode="HTML")
    except Exception as e:
        return await event.reply(f"<i>Oops master, something wrong has happened. \n\nError:</i> <code>{e}</code>", parse_mode="HTML")
