import discord
import os
import src.gaslight as gaslight
import src.help as help
import src.magic as magic
from dotenv import load_dotenv

load_dotenv()

TOKEN = os.getenv('TOKEN')

intents = discord.Intents.default()
intents.message_content = True

client = discord.Client(intents=intents)

@client.event
async def on_message(message):
    if message.author == client.user:
        return

    if message.content.lower().startswith('!gaslight'):
        await message.channel.send(gaslight.gaslight())

    if message.content.lower().startswith('!card'):
        await message.channel.send(magic.get_card(message.content))

    if message.content.lower().startswith('!set'):
        await message.channel.send(magic.set_lookup(message.content))
    
    if message.content.lower().startswith('!help'):
        await message.channel.send(help.helper_message())

client.run(TOKEN)