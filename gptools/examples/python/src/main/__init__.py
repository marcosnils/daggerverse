import dagger
from dagger import dag, function, object_type
from collections.abc import Coroutine


@object_type
class Python:
    @function
    async def gptools(
        self,
        openai_api_key: dagger.Secret,
        question: str,
    ) -> str:
        """example on how to run a full e2e RAG model across different types of
        documents in a directory"""
        nixPaper = dag.http(
            "https://edolstra.github.io/pubs/nspfssd-lisa2004-final.pdf",
        )
        foxImage = dag.http(
            "https://fsquaredmarketing.com/wp-content/uploads/2024/04/bitter-font.png",
        )
        return await dag.gptools().rag(
            openai_api_key,
            dag.directory()
            .with_file("nix-paper.pdf", nixPaper)
            .with_file("image.png", foxImage),
            question,
        )

    @function
    async def gptools_transcript(
        self,
        url: str = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
    ) -> str:
        """example on how to get a transcript for a video and"""
        # TODO: call transcript once python gen is fixed
        # video = dag.http(url)
        return await dag.gptools().audio(url).contents()

    @function
    async def gptools_transcript__directory(
        self,
        url: str = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
    ) -> dagger.Directory:
        """example on how to get a transcript for a video and
        return it as a *Directory to use in other functions"""
        # TODO: call transcript once python gen is fixed
        # video = dag.http(url)
        return dag.directory().with_file(
            "video-transcript.txt", dag.gptools().audio(url)
        )
